import { animate, τmax } from "../motion.js";

export default function (el) {
  const cars = el.querySelector("g");
  animate(cars, { rotate: [0, 9, 0] }, { duration: τmax });
}
