import { animate, τmax } from "../motion.js";

export default function (el) {
  const feather = el.querySelector("g");
  animate(feather, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
