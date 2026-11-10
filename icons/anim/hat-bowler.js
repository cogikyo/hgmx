import { animate, τmax } from "../motion.js";

export default function (el) {
  const hat = el.querySelector("g");
  animate(hat, { rotate: [0, -14, 8, -4, 0] }, { duration: τmax });
}
