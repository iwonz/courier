class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.7/courier_0.3.7_source.tar.gz"
  sha256 "3261f06d0e332b73e7b174b04e7e75ed2f6cbfb8404ddd96de945f11ea68a14a"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=4e86fe72caa599786a97d0c1aed07f85af74691c
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-08T10:44:24+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
