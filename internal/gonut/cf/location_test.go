package cf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/homeport/pina-golada/pkg/files"

	"github.com/gonvenience/text"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// getTempTestingDirectory returns a randomized path based on the .gonut folder
// but doesn't create the directories
func getTempTestingDirectory() string {
	dir := fmt.Sprintf("%s/.gonut/.testing/%s", HomeDir(), text.RandomString(12))
	return dir
}

var _ = Describe("Sample App Location", func() {
	const (
		gonutMainTestingRepoURL = "https://github.com/git-fixtures/basic/" // Sample project from go-git which is used for their test suites as well
		gonutNoRemoteGitRepoURL = "https://ibm.com/"
	)

	Describe("home directory retrieval", func() {
		Context("loading from environmental variables", func() {
			var (
				dir string
			)

			It("should not return an empty string", func() {
				dir = HomeDir()
				Expect(dir).ToNot(BeEmpty())
			})
		})
	})

	Describe("git cloning and pulling", func() {
		Context("clone/pull from valid remote git repo", func() {
			var (
				path string
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustBeforeEach(func() {
				err = cloneOrPull(path, gonutMainTestingRepoURL)
			})

			JustAfterEach(func() {
				setupError := os.RemoveAll(path)
				if setupError != nil {
					panic(setupError)
				}
			})

			It("should clone and create local directory when it doesn't exist", func() {
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(path)
				Expect(err).ToNot(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeFalse())
			})

			It("should pull if directory already exists", func() {
				Expect(err).ToNot(HaveOccurred())
				_, err = os.Stat(path)
				Expect(err).ToNot(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeFalse())

				err = cloneOrPull(path, gonutMainTestingRepoURL)
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(path)
				Expect(err).ToNot(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeFalse())
			})

		})

		Context("clone from url which doesn't link to a git repo", func() {
			var (
				path string
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustBeforeEach(func() {
				err = cloneOrPull(path, gonutNoRemoteGitRepoURL)
			})

			JustAfterEach(func() {
				setupErr := os.RemoveAll(path)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			It("should throw error", func() {
				Expect(err).To(HaveOccurred())

				_, err = os.Stat(path)
				Expect(err).To(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})

		Context("try to execute pull in invalid git repo", func() {
			var (
				path string
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustBeforeEach(func() {
				setupErr := os.MkdirAll(path, os.ModePerm)
				if setupErr != nil {
					panic(setupErr)
				}
				err = cloneOrPull(path, gonutNoRemoteGitRepoURL)
			})

			JustAfterEach(func() {
				setupErr := os.RemoveAll(path)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			It("should try to pull and throw error", func() {
				Expect(err.Error()).To(Equal("repository does not exist"))
			})
		})

	})

	Describe("integration test - load sample app by url", func() {
		Context("loading from valid local directory", func() {
			var (
				path string
				dir  files.Directory
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustBeforeEach(func() {
				setupErr := os.MkdirAll(path, os.ModePerm)
				if setupErr != nil {
					panic(setupErr)
				}
				setupErr = ioutil.WriteFile(path+"/testing.txt", []byte("Hello Gonut!"), 0644)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			JustAfterEach(func() {
				setupErr := os.RemoveAll(path)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			It("should return an in-memory copy of the directory if url is valid", func() {
				dir, err = LoadSampleAppURL(fmt.Sprintf("file:%s", path), "", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(dir.Files())).To(Equal(1))
			})
		})

		Context("loading from invalid local directory", func() {
			var (
				path string
				err  error
			)

			JustBeforeEach(func() {
				path = getTempTestingDirectory()
			})

			It("should return an error when trying to access non-existing directory", func() {
				_, err = LoadSampleAppURL(fmt.Sprintf("file:%s", path), "", "")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("loading from valid git repo", func() {
			var (
				path string
				dir  files.Directory
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustAfterEach(func() {
				setupErr := os.RemoveAll(path)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			It("should return an in-memory copy of repo directory if url is valid", func() {
				dir, err = LoadSampleAppURL(gonutMainTestingRepoURL, "", path)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(dir.Files())).To(BeNumerically(">", 0))

				fileNames := make([]string, len(dir.Files()))
				for index := range dir.Files() {
					fileNames[index] = dir.Files()[index].AbsolutePath().String()
				}
				Expect(fileNames).To(ContainElement("/.gitignore"))
			})

			It("should use relative path and only return contents of sub-directory", func() {
				dir, err = LoadSampleAppURL(gonutMainTestingRepoURL, "go", path)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(dir.Files())).To(BeNumerically(">", 0))

				fileNames := make([]string, len(dir.Files()))
				for index := range dir.Files() {
					fileNames[index] = dir.Files()[index].AbsolutePath().String()
				}
				Expect(fileNames).ToNot(ContainElement("/.gitignore"))
			})
		})

		Context("loading from invalid git repo URL", func() {
			var (
				path string
				dir  files.Directory
				err  error
			)

			BeforeEach(func() {
				path = getTempTestingDirectory()
			})

			JustAfterEach(func() {
				setupErr := os.RemoveAll(path)
				if setupErr != nil {
					panic(setupErr)
				}
			})

			It("should throw an error", func() {
				dir, err = LoadSampleAppURL(gonutNoRemoteGitRepoURL, "", path)
				Expect(err).To(HaveOccurred())
				Expect(dir).To(BeNil())
			})
		})
	})
})
