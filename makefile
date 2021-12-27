static:
	cd frontend && npm run build
	mkdir -p backend/static
	rsync -Pvr frontend/build/ backend/static/
	cd backend && go build
