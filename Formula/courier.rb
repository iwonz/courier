class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.1/courier_0.3.1_source.tar.gz"
  sha256 "f8197b27e497e67377862fc624fd082708e2a842acb4f389b2f4942a5707b711"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=3b3c6281aa1094079438cb638750c390d87ac5f0
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-20T19:51:35+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
