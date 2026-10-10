class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.13/courier_0.3.13_source.tar.gz"
  sha256 "0841fa7366e4bf659343f7ec3f907fa13dfc81d4a87a4f42aee1fbfc2768cac5"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=3211d45a295aec8dc07f77079b3b7a4e2a918d08
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-10T20:16:36+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
