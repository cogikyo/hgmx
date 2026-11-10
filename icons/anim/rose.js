import { animate, τmax } from "../motion.js";

export default function (el) {
  const rose = el.querySelector("g");
  animate(rose, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
