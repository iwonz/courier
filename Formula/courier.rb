class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.3/courier_0.3.3_source.tar.gz"
  sha256 "dd8a830dc1fd9faf39f5a3a744411c8ef8e23018281beb086284f887cd963583"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=74a11f88921e4e1c381cbc545a6e3987704cc009
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-20T20:33:32+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
