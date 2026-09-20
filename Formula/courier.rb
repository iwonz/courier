class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.2.0/courier_0.2.0_source.tar.gz"
  # Homebrew infers version "0.2.0" from the source archive URL.
  sha256 "7b2afa2e88778d352927d9fe38966fb3a5a267e17c56beafb468608748d65204"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=b49f663d910b1ec518cdf12cc00e33d336e48c26
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-09-20T08:13:27+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
