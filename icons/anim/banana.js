import { animate, τmax } from "../motion.js";

export default function (el) {
  const banana = el.querySelector("g");
  animate(banana, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
