import { animate, τmax } from "../motion.js";

export default function (el) {
  const jug = el.querySelector("g");
  animate(jug, { rotate: [0, -15, 5, 0] }, { duration: τmax });
}
