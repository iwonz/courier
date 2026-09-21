class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.4/courier_0.3.4_source.tar.gz"
  sha256 "86ecff0ce1edb62b81181d73029daf0b4846c9792e7e267965c718aeba926ad9"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=407f214ac457bcf2d9f1553b5ce7d9be98c1738a
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-21T13:44:06+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
