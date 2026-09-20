class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "c554ea03caabb40a818b38cec6bb2a82a8f3512e13740ffd49a2fa3db1d9e56b"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=7635d8312a7e71ddf1b8d319fc6e7b7045134400
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-13T18:46:49+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
