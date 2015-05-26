
%w(linux windows darwin).each do |os|
  %w(amd64 386).each do |arch|
    `GOOS=#{os} GOARCH=#{arch} go build`
  end
end
