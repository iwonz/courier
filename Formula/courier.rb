class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.12/courier_0.3.12_source.tar.gz"
  sha256 "0d4665c9d0de059b6d7d83823176721d846c1c651b5411c40810a66e8a24f90e"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=3efdcd64fbf980abd30d53bd82c5f4800a6180bf
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-09T15:58:14+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
