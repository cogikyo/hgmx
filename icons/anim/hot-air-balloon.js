import { animate, τmax } from "../motion.js";

export default function (el) {
  const balloon = el.querySelector("g");
  animate(balloon, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
