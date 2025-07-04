#!/bin/bash
MAIN_VIDEO="simpleInput/body.mp4"               # 主视频文件：将被切割成片段的原始长视频
SEGMENT_NAMES=("body1_5" "body6_10")  # 片段的自定义名称（用于生成文件名）
SEGMENT_TIMES=("1-6" "7-12")    # 时间区间（秒）：决定从主视频截取哪些片段
HOOK_VIDEOS=("simpleInput/hook1.MOV" "simpleInput/hook2.MP4")

# 还存在的问题：如果hook是mov格式，则需要使用更复杂的处理逻辑，不然会卡顿3s
# ======================================
# 创建输出目录
# ======================================
OUTPUT_DIR="final_output_simple"  
rm -rf "$OUTPUT_DIR" 
mkdir -p "$OUTPUT_DIR/tmp"           # 创建主输出目录（-p确保目录不存在时自动创建）,临时目录存放中间处理文件

# ======================================
# 1. 精确切割主视频片段
# ======================================
echo "===== 步骤1: 精确切割主视频 ====="
body_files=()  # 初始化数组，用于存储生成的片段文件路径
# 循环处理每个预定义的视频片段
for (( i=0; i<${#SEGMENT_NAMES[@]}; i++ )); do
    segment_name=${SEGMENT_NAMES[$i]}  # 获取当前片段名称（如"body1_5"）
    
    # 解析时间区间：将"1-6"分割成开始时间(1)和结束时间(6)
    IFS='-' read -ra times <<< "${SEGMENT_TIMES[$i]}"
    start_time=${times[0]}  # 片段开始时间（秒）
    end_time=${times[1]}    # 片段结束时间（秒）
    
    echo "片段: $segment_name (${start_time}-${end_time}秒)"  # 显示当前处理信息
    
    # 获取主视频总时长（防止结束时间超出视频长度）
    duration=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$MAIN_VIDEO")
    echo "总时长: $duration 秒"
    
    # 检查结束时间是否超过视频长度
    if (( $(echo "$end_time > $duration" | bc -l) )); then
        echo "错误: 片段 '$segment_name' 的结束时间 ${end_time}秒 超出视频长度 $duration 秒"
        exit 1  # 发现错误立即停止脚本
    fi
    
    # 计算片段持续时间（结束时间 - 开始时间）
    segment_duration=$(echo "$end_time - $start_time" | bc)
    
    # 切割视频（保留原始音视频）
    output_file="$OUTPUT_DIR/tmp/$segment_name.mp4"
    ffmpeg -y -v error -i "$MAIN_VIDEO" -ss "$start_time" -t "$segment_duration" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p -c:a aac -b:a 128k -ar 48000 -map 0:v:0 -map 0:a:0? -avoid_negative_ts make_zero "$output_file"
    body_files+=("$output_file")  # 将生成的文件路径添加到数组
    echo "生成片段: $output_file"  # 显示生成信息
done


# ======================================
# 2. 预处理Hook视频（片头视频）
# ======================================
echo "\n===== 步骤2: 预处理Hook视频 ====="
hook_files=()  # 初始化数组，存储处理后的Hook文件路径
# 遍历所有Hook视频文件
for hook in "${HOOK_VIDEOS[@]}"; do # 将 $hook 改为 $hook 以保持一致性
    hook_base=$(basename "$hook")
    hook_name="${hook_base%.*}" # 获取文件名（不含扩展名）
    hook_file="$OUTPUT_DIR/tmp/${hook_name}.mp4" # Hook文件输出路径
    extension="${hook##*.}"     # 获取文件扩展名

    ffmpeg_cmd="ffmpeg -y -v error" # 基础FFmpeg命令，抑制错误以外的输出

    echo "正在处理Hook: $hook → $hook_file"
    
    # 根据扩展名选择处理方式
    if [[ "$extension" =~ ^(MOV|mov)$ ]]; then
        echo "检测到MOV格式，将进行更严格的帧率和时间戳处理..."
        # 增加FFmpeg分析媒体流的深度，有助于处理不规则或损坏的MOV文件
        ffmpeg_cmd+=" -analyzeduration 2147483647 -probesize 2147483647" 
       # 尝试获取MOV视频的原始帧率。如果获取失败，默认使用 30fps。
        original_fps=$(ffprobe -v error -select_streams v:0 -show_entries stream=avg_frame_rate -of default=noprint_wrappers=1:nokey=1 "$hook" | head -n 1) 
        if [[ -z "$original_fps" || "$original_fps" == "0/0" ]]; then
            echo "警告: 无法获取MOV文件的原始帧率，默认使用 30fps。请检查源文件或手动设置一个合适的帧率。"
            original_fps="30"
        else
            # 将分数形式的帧率转换为浮点数，例如 "30000/1001" -> 29.97
            original_fps=$(echo "scale=2; $original_fps" | bc -l)
            echo "检测到原始帧率: ${original_fps}fps"
        fi

        # 输入文件和强制固定帧率 (CFR)
        # -r $original_fps 强制输出固定帧率
        # -vsync cfr 确保视频流同步到恒定帧率
        ffmpeg_cmd+=" -i \"$hook\" -r $original_fps -vsync cfr" 
        
        # 使用滤镜重置视频和音频的时间戳，确保从0开始，处理潜在的不连续问题
        ffmpeg_cmd+=" -filter_complex \"[0:v]setpts=PTS-STARTPTS[v];[0:a]asetpts=PTS-STARTPTS[a]\" -map \"[v]\" -map \"[a]\"" 
        
    else # 对于非MOV文件（如MP4），使用标准的标准化参数
        ffmpeg_cmd+=" -i \"$hook\" -map 0:v:0 -map 0:a:0?" 
    fi
    # 通用编码参数：所有Hook视频都将转换为这些统一的格式
    ffmpeg_cmd+=" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p" # 视频编码 H.264, 预设, 质量, 像素格式
    ffmpeg_cmd+=" -c:a aac -b:a 192k -ar 48000" # 音频编码 AAC, 比特率, 采样率
    ffmpeg_cmd+=" -avoid_negative_ts make_zero" # 避免负时间戳
    ffmpeg_cmd+=" \"$hook_file\"" # 输出文件路径
    
    # 执行构建的 FFmpeg 命令
    eval "$ffmpeg_cmd"

    if [ $? -eq 0 ]; then
        hook_files+=("$hook_file")  # 添加文件路径到数组
        echo "Hook视频处理成功: $hook_file"
    else
        echo "Hook视频 $hook 处理失败。跳过此文件。"
    fi
done

# ======================================
# 3. 为每个Hook生成最终视频
# ======================================
echo "\n===== 步骤3: 处理视频组合 ====="

# 遍历每个Hook视频
for hook_file in "${hook_files[@]}"; do
    hook_name=$(basename "$hook_file" .mp4)  # 获取Hook名称

    # 生成随机顺序：打乱片段顺序实现多样化组合
    body_count=${#body_files[@]}  # 获取片段总数
    shuffled_indices=($(seq 0 $((body_count-1)) | shuf))  # 生成随机索引序列
    
    # 构建序列名称（用于文件名）
    sequence=""
    for idx in "${shuffled_indices[@]}"; do
        seg_name=$(basename "${body_files[$idx]}" .mp4)  # 获取片段名
        sequence="${sequence}${seg_name}_"  # 用下划线连接片段名
    done
    sequence=${sequence%_}  # 移除末尾多余的下划线

    # ======================================
    # 4. 拼接视频（Hook + 随机片段）
    # ======================================
    echo "\n===== 步骤4: 拼接视频 (Hook + 正文片段) ====="
    combined_file="$OUTPUT_DIR/tmp/combined_${hook_name}_${sequence}.mp4"  # 拼接后文件路径
    
    # 创建临时文件列表（FFmpeg拼接所需）
    tmp_list="$OUTPUT_DIR/tmp/concat_list.txt"
    echo "file '$(realpath "$hook_file")'" > "$tmp_list"  # 写入Hook文件绝对路径
    
    # 按随机顺序添加主视频片段
    for idx in "${shuffled_indices[@]}"; do
        echo "file '$(realpath "${body_files[$idx]}")'" >> "$tmp_list"
    done
    
    # 使用FFmpeg拼接视频
    #ffmpeg -y -f concat -safe 0 -i "$tmp_list" -fflags +genpts -avoid_negative_ts make_zero -fps_mode passthrough -filter_complex "[0:v]setpts=PTS-STARTPTS[v];[0:a]asetpts=PTS-STARTPTS[a]" -map "[v]" -map "[a]" -bsf:a aac_adtstoasc "$combined_file"
    #ffmpeg -y -f concat -safe 0 -i "$tmp_list" -c copy "$combined_file"
    # 使用FFmpeg拼接并完全重编码视频和音频
    # -f concat -safe 0 -i "$tmp_list"：使用concat demuxer读取列表文件
    # -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p：视频重编码参数
    # -g 30 -keyint_min 30：强制每秒至少一个关键帧，提高兼容性
    # -c:a aac -b:a 192k -ar 48000：音频重编码参数
    # -avoid_negative_ts make_zero：确保时间戳为正
    # -shortest：确保输出时长与最短流一致
    # -fflags +genpts：有时能帮助生成准确的时间戳（尽管重编码是主要手段）
    # -vsync cfr：在合并时强制固定帧率
    ffmpeg -y -v error -f concat -safe 0 -i "$tmp_list" \
           -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p -g 30 -keyint_min 30 \
           -c:a aac -b:a 192k -ar 48000 \
           -fflags +genpts -vsync cfr \
           -avoid_negative_ts make_zero -shortest \
           "$combined_file"

    # 查看拼接后视频的信息
    echo "\n拼接后视频信息："
    ffprobe -v error -show_entries format=duration,format=size -of default=noprint_wrappers=1:nokey=1 "$combined_file"

    echo "拼接成功: $combined_file"
    #rm "$tmp_list"  # 删除临时文件
done
