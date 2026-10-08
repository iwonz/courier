class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.8/courier_0.3.8_source.tar.gz"
  sha256 "1be959ca4a10e63e56abcab8c022f0508461e57a44d2ead4befedc3c02ee8aff"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=8ec73ba3ee2d1df4d08b52bbedf15a8140824da0
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-08T14:03:07+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
