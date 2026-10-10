class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.16/courier_0.3.16_source.tar.gz"
  sha256 "6a3b552fa9e2467158992720fa8a62ec68e05792fa0cc2eb570c610f286540c8"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=482665080da12b1af126683f6d070a1e2fd526b9
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-10T21:00:04+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
