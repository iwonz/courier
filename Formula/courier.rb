class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.2/courier_0.3.2_source.tar.gz"
  sha256 "7efdb3a7b6f591f7912c3a559b8747609f9e3969f4c952f175300cafdf7f1f7e"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=60d92a8a7cf8b0ca50c04d91b9a25d3c19ffa3b6
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-20T20:13:48+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
