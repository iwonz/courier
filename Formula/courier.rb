class Courier < Formula
  desc "Safely transfer files and directories between local and SSH endpoints"
  homepage "https://github.com/iwonz/courier"
  url "https://github.com/iwonz/courier/releases/download/v0.3.6/courier_0.3.6_source.tar.gz"
  sha256 "3007be3bf34506bf57fe67b91f78c88413a0ce9e19f0c10b549821d258337aca"
  license "MIT"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X github.com/iwonz/courier/internal/buildinfo.Version=#{version}
      -X github.com/iwonz/courier/internal/buildinfo.Commit=8227ccfbab46f6c49b24ab9e11099b2944cc57ef
      -X github.com/iwonz/courier/internal/buildinfo.Date=2026-10-08T09:56:01+03:00
    ]
    system "go", "build", "-trimpath", "-ldflags", ldflags.join(" "), "-o", bin/"courier", "./cmd/courier"
  end

  test do
    assert_match "courier #{version}", shell_output("#{bin}/courier version")
  end
end
