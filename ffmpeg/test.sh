#!/usr/bin/env bash
set -euo pipefail

MAIN_VIDEO="input/body.mp4"
SEGMENT_NAMES=(body1_5 body6_10)
SEGMENT_TIMES=(1-6 7-12)
HOOK_VIDEOS=(input/hook2.mp4)

OUTPUT_DIR="final_output_simple"
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR/tmp"

# 1. 切片
echo "=== 步骤1: 切片 ==="
body_files=()
for i in "${!SEGMENT_NAMES[@]}"; do
  name=${SEGMENT_NAMES[i]}
  IFS='-' read -r start end <<< "${SEGMENT_TIMES[i]}"
  dur=$(echo "$end - $start" | bc -l)
  out="$OUTPUT_DIR/tmp/${name}.mp4"
  ffmpeg -y -v error \
    -ss "$start" -i "$MAIN_VIDEO" \
    -t "$dur" \
    -c:v libx264 -preset fast -crf 23 \
    -c:a aac -b:a 128k \
    "$out"
  echo "Slice -> $out"
  body_files+=("$out")
done

# 2. 预处理 Hook（全部转到 H.264/AAC）
echo "=== 步骤2: 预处理 Hook ==="
hook_files=()
for hook in "${HOOK_VIDEOS[@]}"; do
  base=$(basename "$hook" | sed 's/\..*$//')
  out="$OUTPUT_DIR/tmp/${base}.mp4"
  ffmpeg -y -v error -i "$hook" \
    -c:v libx264 -preset fast -crf 23 \
    -c:a aac -b:a 128k \
    "$out"
  hook_files+=("$out")
  echo "Hook -> $out"
done

# 3. 随机打乱函数
shuffle_array() {
  local arr_name=\$1
  local -a arr=(\${!arr_name[@]})
  for ((i=\${#arr[@]}-1; i>0; i--)); do
    j=\$(( RANDOM % (i+1) ))
    tmp=\${arr[i]}; arr[i]=\${arr[j]}; arr[j]=\$tmp
  done
  \$arr_name=(\${arr[@]})
}

# 4. 拼接
echo "=== 步骤3: 拼接 ==="
for hook in "${hook_files[@]}"; do
  name=$(basename "$hook" .mp4)
  # 打乱副本
  to_concat=("${body_files[@]}")
  shuffle_array to_concat

  # 写列表
  list="$OUTPUT_DIR/tmp/${name}_list.txt"
  : > "$list"
  echo "file '$hook'" >> "$list"
  for seg in "${to_concat[@]}"; do
    echo "file '$seg'" >> "$list"
  done

  out="$OUTPUT_DIR/${name}_final.mp4"
  ffmpeg -y -v error \
    -f concat -safe 0 -i "$list" \
    -c copy \
    "$out"
  echo "生成 -> $out"
  rm -f "$list"
done