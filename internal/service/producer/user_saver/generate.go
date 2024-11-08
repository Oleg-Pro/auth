package user_saver

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../../bin/minimock -i UserSaverProducer -o ./mocks/ -s "_minimock.go"
