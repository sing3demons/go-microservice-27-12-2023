#

```cmd
kubectl -n ms-service get po
kubectl -n ms-service logs -f po/
```

```k8s
kubectl apply -f ./00-namespace.yml
```

```mongodb
cd mongodb
kubectl apply -f ./01-pv-pvc.yml 
kubectl apply -f ./02-statefulset.yml
kubectl apply -f ./03-service.yml

kubectl exec -it mongo-0 -n ms-service -- mongosh --eval 'rs.initiate({_id: "rs0", members: [{ _id: 0, host: "mongo-0.mongo:27017" },{ _id: 1, host: "mongo-1.mongo:27017" },{ _id: 2, host: "mongo-2.mongo:27017" }]})'
```

```products
cd products
kubectl apply -f ./01-deployment.yml
kubectl apply -f ./02-service.yml
```
