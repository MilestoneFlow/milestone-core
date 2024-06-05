	(
	  chmod +x ./milestone_core
	  export $(grep -v '^#' .env | xargs) && go build && ./milestone_core
	)