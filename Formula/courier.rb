class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.5/courier_0.3.5_source.tar.gz"
  sha256 "32812fa3567541b3755bfb1d4900fa3cc977ebfe7d26eb753d8cce543db74037"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=f84c678520575704388e4936c3e3ac8c1dd41ae4
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-21T17:34:22+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
