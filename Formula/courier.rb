class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.10/courier_0.3.10_source.tar.gz"
  sha256 "0157377f08f37af4f77ce13e6b9cfeffa664001a8fbfe4e2b7fdd34a6b5548f6"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=0aad5c52eab511d463ef0ca21c927975a032db4b
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-09T11:38:34+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
