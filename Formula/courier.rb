class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.19/courier_0.3.19_source.tar.gz"
  sha256 "02793a52f98f4b7bee2956dc3ee50b7251347f3112890b1a7bf5ec19aacde9a7"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=7e6fa87894c24adcae62d51f08922069d7c6ff7a
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-11T09:54:37+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
