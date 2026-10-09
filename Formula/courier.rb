class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.11/courier_0.3.11_source.tar.gz"
  sha256 "28d63ddb7ce63dc702fc4a2089ae598d3484d05f0eb2ef76c4c678de55d71c36"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=384fd60915327aaf5e0713a9dc42ee4de749fb59
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-09T15:20:39+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
