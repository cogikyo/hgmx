import { animate, τmax } from "../motion.js";

export default function (el) {
  const flower = el.querySelector("g");
  animate(flower, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
