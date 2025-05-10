# Kubernetes Simple Controller

A simple Kubernetes controller written in Go to watch for Pod events (Add, Update, Delete) and process them efficiently using a workqueue.

## 📌 **Project Structure**

```
├── main.go
├── go.mod
├── go.sum
└── README.md
```

* **main.go**: Main controller logic to watch and process Pod events.
* **go.mod**: Go module dependencies.
* **go.sum**: Dependency checksums.
* **README.md**: Project documentation.

---

## 🚀 **Getting Started**

### Prerequisites

* Kubernetes Cluster (Minikube, Kind, EKS, etc.)
* kubeconfig file configured (`~/.kube/config`)
* Go (1.18 or later)

### Installation

Clone the repository:

```sh
git clone <repo-url>
cd Kubernetes-Simple-Controller
```

Initialize Go modules and install dependencies:

```sh
go mod tidy
```

---

## 🔎 **Usage**

Set your kubeconfig path:

```sh
export KUBECONFIG=~/.kube/config
```

Run the controller:

```sh
go run main.go
```

To build the executable:

```sh
go build -o k8s-controller
./k8s-controller
```

---

## 📋 **Test the Controller**

Create a test pod:

```sh
kubectl run nginx-pod --image=nginx --restart=Never
```

You should see logs indicating the pod has been added.

Delete the test pod:

```sh
kubectl delete pod nginx-pod
```

---

## 🤝 **Contributing**

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`)
3. Make your changes and commit (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Create a Pull Request

---

## 🛡️ **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact**

For any questions or suggestions, feel free to reach out or open an issue.
