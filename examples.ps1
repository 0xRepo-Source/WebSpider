# WebSpider - Example Usage Scenarios (PowerShell)

Write-Host "WebSpider - Example Usage Scenarios" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green

# Example 1: Discover directory structure with rate limiting
Write-Host ""
Write-Host "1. Discover PDF files from a hypothetical university site:" -ForegroundColor Yellow
Write-Host ".\webspider.exe -url `"https://university.edu/publications/`" ``" -ForegroundColor Cyan
Write-Host "  -discover-only ``" -ForegroundColor Cyan
Write-Host "  -rate 0.5 ``" -ForegroundColor Cyan
Write-Host "  -accept `"\.(pdf|doc|docx)$`" ``" -ForegroundColor Cyan
Write-Host "  -depth 4 ``" -ForegroundColor Cyan
Write-Host "  -save-list `"academic-papers.txt`" ``" -ForegroundColor Cyan
Write-Host "  -verbose" -ForegroundColor Cyan

Write-Host ""
Write-Host "2. Download software packages with conservative rate limiting:" -ForegroundColor Yellow
Write-Host ".\webspider.exe -url `"https://releases.example.com/`" ``" -ForegroundColor Cyan
Write-Host "  -discover-only ``" -ForegroundColor Cyan
Write-Host "  -rate 0.25 ``" -ForegroundColor Cyan
Write-Host "  -accept `"\.(tar\.gz|zip|deb|rpm)$`" ``" -ForegroundColor Cyan
Write-Host "  -depth 2 ``" -ForegroundColor Cyan
Write-Host "  -save-list `"packages.txt`"" -ForegroundColor Cyan
Write-Host ""
Write-Host "# After editing packages.txt to select desired files:" -ForegroundColor Gray
Write-Host ".\webspider.exe -urls `"packages.txt`" -rate 0.5 -output `".\software`"" -ForegroundColor Cyan

Write-Host ""
Write-Host "3. Mirror documentation while excluding assets:" -ForegroundColor Yellow
Write-Host ".\webspider.exe -url `"https://docs.example.com/`" ``" -ForegroundColor Cyan
Write-Host "  -rate 1.5 ``" -ForegroundColor Cyan
Write-Host "  -reject `"\.(jpg|jpeg|png|gif|svg|css|js)$`" ``" -ForegroundColor Cyan
Write-Host "  -accept `"\.(html|htm|pdf|txt|md)$`" ``" -ForegroundColor Cyan
Write-Host "  -depth 5 ``" -ForegroundColor Cyan
Write-Host "  -output `".\docs-mirror`"" -ForegroundColor Cyan

Write-Host ""
Write-Host "4. Very conservative crawling for rate-limited servers:" -ForegroundColor Yellow
Write-Host ".\webspider.exe -url `"https://sensitive-server.com/files/`" ``" -ForegroundColor Cyan
Write-Host "  -discover-only ``" -ForegroundColor Cyan
Write-Host "  -rate 0.1 ``" -ForegroundColor Cyan
Write-Host "  -timeout 60s ``" -ForegroundColor Cyan
Write-Host "  -depth 3 ``" -ForegroundColor Cyan
Write-Host "  -user-agent `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`" ``" -ForegroundColor Cyan
Write-Host "  -verbose" -ForegroundColor Cyan

Write-Host ""
Write-Host "Remember:" -ForegroundColor Magenta
Write-Host "- Start with discovery-only mode (-discover-only)" -ForegroundColor White
Write-Host "- Use low rate limits (0.1-1.0 req/sec) for respectful crawling" -ForegroundColor White
Write-Host "- Edit the generated URL list before downloading" -ForegroundColor White
Write-Host "- Monitor with -verbose flag for debugging" -ForegroundColor White

Write-Host ""
Write-Host "Quick Test (safe example):" -ForegroundColor Green
Write-Host ".\webspider.exe -url `"https://httpbin.org/`" -discover-only -rate 1.0 -depth 1 -verbose" -ForegroundColor Yellow