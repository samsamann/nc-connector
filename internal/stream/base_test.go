package stream

type testProducer struct {
}

/*type outerItem struct {
	SyncItem
}

func (p testProducer) In() <-chan SyncItem {
	c := make(chan SyncItem)
	go func() {
		defer close(c)
		for i := 0; i < 10; i++ {
			c <- outerItem{}
			time.Sleep(time.Second)
		}
	}()
	return c
}

func TestProducer(t *testing.T) {
	c := (testProducer{}).In()
	for v := range c {
		t.Logf("%v", v)
	}
}
*/
