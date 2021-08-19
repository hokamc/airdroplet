for filename in githook/*; do
  filename=$(basename ${filename})
	cp "githook/$filename" ".git/hooks/$filename"; chmod +x ".git/hooks/$filename" &
done