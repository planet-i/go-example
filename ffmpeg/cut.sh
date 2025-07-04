#!/bin/bash
#
# 视频转换与裁剪脚本
# 功能：将 MOV/MP4 视频转换为标准 MP4 格式，并精确截取指定片段。
# 目标：确保输出的视频文件格式高度统一和标准化，方便后续的拼接操作，避免兼容性问题。
#
# 用法：./video_convert_cut.sh [输出目录] [起始时间(可选)] [保留时间(可选)]
# 示例：sh cut.sh hook 00:00:00 00:00:03  (所有预定义视频将裁剪3秒，从开头)
# 示例：sh cut.sh body 00:00:03 00:00:05  (所有预定义视频将从3秒处裁剪5秒)

# 示例：假设脚本和 'cutInput' 文件夹在同一个父目录下
HOOK_VIDEOS=(
    "cutInput/hook/hook1.MOV"
    "cutInput/hook/hook2.MOV"
    "cutInput/hook/hook3.MOV"
)

MAIN_VIDEOS=(
    "cutInput/body/body1.mp4"
    "cutInput/body/body2.mp4" 
    "cutInput/body/body3.mp4"
)

# --- 参数检查 ---
if [ $# -lt 1 ]; then # 只需要输出目录，起始时间和持续时间是可选的
    echo "用法: $0 <输出目录> [起始时间，默认00:00:00] [保留时间，默认00:00:05]"
    echo "提示: 确保输出的视频文件格式统一，方便后续拼接，减少问题。"
    exit 1
fi

VIDEO_TYPE="$1" # 第一个参数：'hook' 或 'body'
# 输出目录将由第一个参数提供
OUTPUT_DIR="cut_output/$1"         # 输出目录：生成的所有文件将保存在此文件夹
rm -rf "$OUTPUT_DIR" 
mkdir -p "$OUTPUT_DIR" 
START_TIME="${2:-00:00:00}"  # 默认从 00:00:00 (零秒) 开始裁剪
DURATION="${3:-00:00:05}"    # 默认裁剪 00:00:05 (5秒) 时长的视频片段

# 将所有要处理的视频路径合并到一个数组中，方便统一遍历
# ALL_VIDEOS=("${HOOK_VIDEOS[@]}" "${MAIN_VIDEOS[@]}")
case "$VIDEO_TYPE" in
    hook)
        ALL_VIDEOS=("${HOOK_VIDEOS[@]}")
        ;;
    body)
        ALL_VIDEOS=("${MAIN_VIDEOS[@]}")
        ;;
    *)
        echo "错误: 无效的视频类型 '$VIDEO_TYPE'。请使用 'hook' 或 'body'。"
        exit 1
        ;;
esac

# --- 创建输出目录 ---
mkdir -p "$OUTPUT_DIR"

echo "====================================="
echo "  视频裁剪与标准化处理开始"
echo "  总共将处理 ${#ALL_VIDEOS[@]} 个视频文件。"
echo "  输出目录:   $OUTPUT_DIR"
echo "  裁剪起始:   $START_TIME"
echo "  裁剪时长:   $DURATION"
echo "====================================="

# --- 遍历预定义的所有视频文件 ---
# 循环处理 ALL_VIDEOS 数组中的每一个视频路径
for input_file_raw in "${ALL_VIDEOS[@]}"; do
    # **关键：将 input_file_raw 转换为绝对路径**
    # 这样传递给 FFmpeg 的路径始终是完整的，避免 FFmpeg 内部的路径解析错误
    input_file="$(realpath "$input_file_raw")"

    # 检查文件是否存在且可读
    if [ ! -f "$input_file" ]; then
        echo "错误: 视频文件 '$input_file_raw' 不存在或无法访问，跳过处理。"
        continue # 跳过当前文件，处理下一个
    fi

    # 生成输出文件名：保留原始文件名，但扩展名统一改为 .mp4
    filename=$(basename -- "$input_file")
    filename_noext="${filename%.*}" # 获取不带扩展名的文件名
    output_file="$OUTPUT_DIR/${filename_noext}.mp4"

    echo "---"
    echo "正在处理文件: $filename"
    echo "完整输入路径: $input_file" # 打印完整的输入路径，方便调试
    echo "输出文件路径: $output_file" # 打印完整的输出路径，方便调试


    extension="${filename##*.}"      # 获取文件扩展名 (如 mov, mp4)
    ffmpeg_cmd="ffmpeg -y -v error"  # 构建 FFmpeg 基础命令：
                                     # -y: 覆盖同名输出文件而不询问
                                     # -v error: 只显示错误信息，保持输出简洁

    # --- 针对 MOV 文件的特殊处理逻辑 ---
    # MOV 文件（特别是来自手机的）常有可变帧率 (VFR) 和复杂的时间戳，需要特殊处理以确保裁剪和拼接的准确性。
    if [[ "$extension" =~ ^(MOV|mov)$ ]]; then
        echo "检测到 MOV 格式，将进行更严格的帧率和时间戳处理..."
        # 增加 FFmpeg 分析媒体流的深度和探测大小，有助于处理不规则或损坏的 MOV 文件。
        ffmpeg_cmd+=" -analyzeduration 2147483647 -probesize 2147483647" 

        # 尝试获取 MOV 视频的原始帧率，用于强制固定帧率 (CFR)。
        # 如果获取失败或帧率为 0/0 (未知)，则默认使用 30fps。
        original_fps=$(ffprobe -v error -select_streams v:0 -show_entries stream=avg_frame_rate -of default=noprint_wrappers=1:nokey=1 "$input_file" | head -n 1)
        if [[ -z "$original_fps" || "$original_fps" == "0/0" ]]; then
            echo "警告: 无法获取 MOV 文件 '$filename' 的原始帧率，默认使用 30fps。请检查源文件或手动设置一个合适的帧率。"
            original_fps="30"
        else
            # 将分数形式的帧率 (如 "30000/1001") 转换为浮点数 (如 29.97)。
            original_fps=$(echo "scale=2; $original_fps" | bc -l)
            echo "检测到原始帧率: ${original_fps}fps"
        fi

        # **关键步骤：使用滤镜重置 MOV 视频和音频的时间戳并强制固定帧率。**
        # 这一步在裁剪之前执行，确保时间线是线性的，裁剪更准确。
        ffmpeg_cmd+=" -i \"$input_file\" -filter_complex \"[0:v]setpts=PTS-STARTPTS,fps=${original_fps}[v];[0:a]asetpts=PTS-STARTPTS[a]\""
        ffmpeg_cmd+=" -map \"[v]\" -map \"[a]\"" # 映射滤镜处理后的视频和音频流
        # 裁剪参数：在时间戳重置和帧率固定后应用，确保裁剪精度。
        ffmpeg_cmd+=" -ss \"$START_TIME\" -t \"$DURATION\""
    # --- 针对非 MOV 文件 (如 MP4) 的通用处理逻辑 ---
    else
        # 对于 MP4 等文件，直接指定输入文件和映射流。
        ffmpeg_cmd+=" -i \"$input_file\" -map 0:v:0 -map 0:a:0?" # 映射第一个视频流和可选的第一个音频流
        # 裁剪参数：对于 MP4，将 -ss 和 -t 放在 -i 之后通常也能保证精度。
        ffmpeg_cmd+=" -ss \"$START_TIME\" -t \"$DURATION\""
        # 注意：这里不需要 setpts 滤镜，因为 MP4 通常时间戳相对规整，且最终输出会重编码。
    fi

    # --- 通用标准化编码参数 (适用于所有输出文件) ---
    # 这些参数与主拼接脚本中的参数保持一致，以确保最终输出视频的兼容性和稳定性。
    ffmpeg_cmd+=" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p"
    ffmpeg_cmd+=" -c:a aac -b:a 192k -ar 48000"
    ffmpeg_cmd+=" -avoid_negative_ts make_zero"
    ffmpeg_cmd+=" \"$output_file\"" # 输出文件路径

    # --- 执行构建的 FFmpeg 命令 ---
    echo "执行命令: $ffmpeg_cmd"
    eval "$ffmpeg_cmd"

    # --- 检查 FFmpeg 命令执行结果 ---
    if [ $? -eq 0 ]; then # $? 是上一条命令的退出状态码，0 表示成功
        echo "成功生成: $(basename "$output_file")"
    else
        echo "处理文件 '$filename' 失败。请检查错误信息。"
    fi
done

echo "====================================="
echo "  所有文件处理完成！输出目录: $OUTPUT_DIR"
echo "====================================="