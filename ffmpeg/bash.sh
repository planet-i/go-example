#!/bin/bash
#set -x  # 开启调试模式（可选）：取消注释后，执行时会显示每条命令及其参数，用于排查错误

# ====================================== 用剪好的片段组合
# BODY_SEGMENT=("cut/body1.mp4" "cut/body2.mp4" "cut/body3.mp4")
# BODY_SEGMENT=("cut/body1.mp4")
# HOOK_VIDEOS=("cut/hook1.mp4" "cut/hook2.mp4") 
# ====================================== 传一个body去裁剪
HOOK_VIDEOS=("input/hook1.MOV") # Hook视频列表：开头插入的短视频（如片头广告）
MAIN_VIDEO="input/body.mp4"               # 主视频文件：将被切割成片段的原始长视频
# 定义命名切割片段（片段名和时间区间）
SEGMENT_NAMES=("body1_5" "body6_10" "body11_15")  # 片段的自定义名称（用于生成文件名）
SEGMENT_TIMES=("1-6" "7-12" "13-18")              # 时间区间（秒）：决定从主视频截取哪些片段

BG_MUSICS=("input/bgm.mp3")             # 背景音乐列表（可为空）：添加到最终视频的背景音乐
VOICE_OVERS=("input/VO1.mp3" "input/VO2.mp3") # 配音音乐列表（可为空）：添加到最终视频的画外音

# ======================================
# 创建输出目录
# ======================================
OUTPUT_DIR="final_output"         # 输出目录：生成的所有文件将保存在此文件夹
rm -rf "$OUTPUT_DIR" 
mkdir -p "$OUTPUT_DIR/tmp"        # 创建主输出目录（-p确保目录不存在时自动创建）,临时目录存放中间处理文件

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
    duration=$(ffprobe -v quiet -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$MAIN_VIDEO")
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
    #                                 平衡编码速度和压缩效率、恒定速率因子、通用像素格式           H.264 编码，veryfast 预设，23 CRF         音频编码 设置统一的音频比特率和采样率  确保我们映射输入中的第一个视频流，并可选地映射第一个音频流  保持时间戳为正
    ffmpeg -y -v quiet -i "$MAIN_VIDEO" -ss "$start_time" -t "$segment_duration" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p -c:a aac -b:a 192k -ar 48000 -map 0:v:0 -map 0:a:0? -avoid_negative_ts make_zero "$output_file"
    
    body_files+=("$output_file")  # 将生成的文件路径添加到数组
    echo "生成片段: $output_file"  # 显示生成信息
done
# 如果直接提供的是片段，就把片段添加进去
for body in "${BODY_SEGMENT[@]}"; do
    body_base=$(basename "$body")
    body_name="${body_base%.*}" # 获取文件名（不含扩展名）
    body_file="$OUTPUT_DIR/tmp/$body_name.mp4" # Hook文件输出路径
    #cp "$body" "$body_file"
    ffmpeg -y -v quiet -i "$body" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p -c:a aac -b:a 192k -ar 48000 -map 0:v:0 -map 0:a:0? -avoid_negative_ts make_zero "$body_file"
    body_files+=("$body_file")  # 添加文件路径到数组
    echo "准备body: $body_file"  # 显示处理信息
done

# ======================================
# 2. 预处理Hook视频（片头视频）
# ======================================
echo "\n===== 步骤2: 预处理Hook视频 ====="
hook_files=()  # 初始化数组，存储处理后的Hook文件路径
# 遍历所有Hook视频文件
for hook in "${HOOK_VIDEOS[@]}"; do
    hook_base=$(basename "$hook")
    hook_name="${hook_base%.*}" # 获取文件名（不含扩展名）
    hook_file="$OUTPUT_DIR/tmp/$hook_name.mp4" # Hook文件输出路径
    extension="${hook##*.}"     # 获取文件扩展名
    ffmpeg_cmd="ffmpeg -y -v quiet" # 基础FFmpeg命令，抑制错误以外的输出
    # 根据扩展名选择处理方式
    if [[ "$extension" =~ ^(MOV|mov)$ ]]; then
        echo "检测到MOV格式，将进行更严格的帧率和时间戳处理...: $hook → $hook_file"
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
        # MOV转MP4并统一参数
        #ffmpeg -y -v quiet -i "$hook" -c:v libx264 -crf 23 -preset fast -c:a aac -ar 44100 -avoid_negative_ts make_zero "$hook_file"
    else # 对于非MOV文件（如MP4），使用标准的标准化参数
        ffmpeg_cmd+=" -i \"$hook\" -map 0:v:0 -map 0:a:0?" # 映射第一个视频流和可选的第一个音频流
        #ffmpeg -y -v quiet -i "$hook" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p -c:a aac -b:a 192k -ar 48000 -map 0:v:0 -map 0:a:0? -avoid_negative_ts make_zero "$hook_file"
    fi
    # 通用编码参数：所有Hook视频都将转换为这些统一的格式
    ffmpeg_cmd+=" -c:v libx264 -preset veryfast -crf 23 -pix_fmt yuv420p" # 视频编码 H.264, 预设, 质量, 像素格式
    ffmpeg_cmd+=" -c:a aac -b:a 192k -ar 48000" # 音频编码 AAC, 比特率, 采样率
    ffmpeg_cmd+=" -avoid_negative_ts make_zero" # 避免负时间戳
    ffmpeg_cmd+=" \"$hook_file\"" # 输出文件路径
    # 执行构建的 FFmpeg 命令
    eval "$ffmpeg_cmd"

    hook_files+=("$hook_file")  # 添加文件路径到数组
    echo "准备Hook: $hook_file"  # 显示处理信息
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
    # ffmpeg -y -v quiet -v quiet -f concat -safe 0 -i "$tmp_list" -fflags +genpts -avoid_negative_ts make_zero -fps_mode passthrough -c copy -bsf:a aac_adtstoasc "$combined_file"
    # ffmpeg -y -v quiet -f concat -safe 0 -i "$tmp_list" -c copy -avoid_negative_ts make_zero "$combined_file"
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
    echo "拼接成功: $combined_file"
    #rm "$tmp_list"  # 删除临时文件

    # ======================================
    # 5. 音频处理与最终输出（分4种情况处理）
    # ======================================
    echo "\n===== 步骤5: 音频处理与最终输出 ====="
    # 检查是否有背景音乐/配音
    has_bg=$([ ${#BG_MUSICS[@]} -gt 0 ] && echo 1 || echo 0)
    has_vo=$([ ${#VOICE_OVERS[@]} -gt 0 ] && echo 1 || echo 0)
    
    # 情况1：无背景音乐和配音
    if (( has_bg == 0 && has_vo == 0 )); then
        output_file="$OUTPUT_DIR/${hook_name}_${sequence}.mp4"
        #cp "$combined_file" "$output_file"  # 直接复制文件
        ffmpeg -y -v quiet -i "$concat_video_file" -c copy "$output_file"
        echo "生成带原声版本: $output_file"
    # 情况2：有背景音乐无配音
    elif (( has_bg == 1 && has_vo == 0 )); then
        for bg in "${BG_MUSICS[@]}"; do
            bg_name=$(basename "$bg" .mp3)
            output_file="$OUTPUT_DIR/${hook_name}_${sequence}_${bg_name}.mp4"
            # 音频处理：降低背景音乐音量（避免盖过人声）
            # ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -filter_complex "[1:a]volume=0.5[bg];[bg][0:a]amerge=inputs=2[a]" -map 0:v -map "[a]" -c:v copy -c:a aac -shortest "$output_file"
            # ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -filter_complex "[1:a]volume=0.2[bgm];[0:a:0][bgm]amerge=inputs=2[aout]" -map 0:v:0 -map "[aout]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file" 有视频原声
            
            # 音频处理：降低背景音乐音量，并替换原始视频音频
            # [0:v:0] 映射视频流
            # [1:a] 映射背景音乐，音量调整
            # -map 0:v:0 只映射视频流，不映射原始音频
            # -map "[bgm]" 映射处理过的背景音乐
            ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -filter_complex "[1:a]volume=0.2[bgm]" -map 0:v:0 -map "[bgm]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file"
            
            echo "生成背景音乐版本: $output_file"
            echo "  Hook: $hook_name | 序列: $sequence | 背景音乐: $bg_name"
        done
    # 情况3：有配音无背景音乐
    elif (( has_bg == 0 && has_vo == 1 )); then
        for vo in "${VOICE_OVERS[@]}"; do
            vo_name=$(basename "$vo" .mp3)
            output_file="$OUTPUT_DIR/${hook_name}_${sequence}_${vo_name}.mp4"
            
            # 音频处理：提高配音音量
            # -filter_complex "[1:a]volume=0.2[bgm];[0:a:0][bgm]amerge=inputs=2[aout]"
            # [1:a] 表示第二个输入文件（背景音乐）的音频流
            # volume=0.2 降低背景音乐音量
            # [bgm] 给处理后的背景音乐流命名
            # [0:a:0] 表示第一个输入文件（拼接视频）的第一个音频流
            # amerge=inputs=2 将两个音频流合并
            # [aout] 给合并后的音频流命名
            # -map 0:v:0 映射第一个输入文件的视频流
            # -map "[aout]" 映射合并后的音频流
            # -c:v copy 复制视频流
            # -c:a aac -b:a 192k 重新编码音频为 AAC 192kbps
            # -shortest 确保输出时长与最短的流一致，防止音画不同步
            # ffmpeg -y -v quiet -i "$combined_file" -i "$vo" -filter_complex "[1:a]volume=0.8[vo];[vo][0:a]amerge=inputs=2[a]" -map 0:v -map "[a]" -c:v copy -c:a aac -shortest "$output_file"
            # ffmpeg -y -v quiet -i "$combined_file" -i "$vo" -filter_complex "[1:a]volume=1.0[vo];[0:a:0][vo]amerge=inputs=2[aout]" -map 0:v:0 -map "[aout]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file"
            
            # 音频处理：提高配音音量，并替换原始视频音频
            # [0:v:0] 映射视频流
            # [1:a] 映射配音，音量调整
            # -map 0:v:0 只映射视频流，不映射原始音频
            # -map "[vo]" 映射处理过的配音
            ffmpeg -y -v quiet -i "$combined_file" -i "$vo" -filter_complex "[1:a]volume=1.0[vo]" -map 0:v:0 -map "[vo]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file"
           
            echo "生成配音版本: $output_file"
            echo "  Hook: $hook_name | 序列: $sequence | 配音: $vo_name"
        done
    # 情况4：同时有背景音乐和配音
    else
        for bg in "${BG_MUSICS[@]}"; do
            for vo in "${VOICE_OVERS[@]}"; do
                bg_name=$(basename "$bg" .mp3)
                vo_name=$(basename "$vo" .mp3)
                output_file="$OUTPUT_DIR/${hook_name}_${sequence}_${bg_name}_${vo_name}.mp4"
                
                # 音频处理：同时调整背景音乐和配音音量
                # [1:a] 是背景音乐，[2:a] 是配音
                # amerge=inputs=3 合并三个音频流：原始视频音频、背景音乐、配音
                # ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -i "$vo" -filter_complex "[1:a]volume=0.2[bg];[2:a]volume=0.8[vo];[bg][vo]amerge=inputs=2[a]" -map 0:v -map "[a]" -c:v copy -c:a aac -shortest "$output_file"
                # ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -i "$vo" -filter_complex "[1:a]volume=0.2[bgm];[2:a]volume=1.0[vo];[0:a:0][bgm][vo]amerge=inputs=3[aout]" -map 0:v:0 -map "[aout]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file"
                
                # 音频处理：同时调整背景音乐和配音音量，替换原始视频音频
                # [0:v:0] 映射视频流
                # [1:a] 是背景音乐，音量调整
                # [2:a] 是配音，音量调整
                # [bgm][vo]amerge=inputs=2[aout] 合并背景音乐和配音
                # -map 0:v:0 只映射视频流，不映射原始音频
                # -map "[aout]" 映射合并后的背景音乐和配音
                ffmpeg -y -v quiet -i "$combined_file" -i "$bg" -i "$vo" -filter_complex "[1:a]volume=0.2[bgm];[2:a]volume=1.0[vo];[bgm][vo]amerge=inputs=2[aout]" -map 0:v:0 -map "[aout]" -c:v copy -c:a aac -b:a 192k -shortest "$output_file"

                echo "生成完整音频版本: $output_file"
                echo "  Hook: $hook_name | 序列: $sequence | 背景音乐: $bg_name | 配音: $vo_name"
            done
        done
    fi
done

# ======================================
# 完成处理
# ======================================
echo "\n===== 处理完成 ====="
echo "所有输出文件保存在: $OUTPUT_DIR"