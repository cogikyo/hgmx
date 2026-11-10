import { animate, τmax } from "../motion.js";

export default function (el) {
  const roller = el.querySelector("g");
  animate(roller, { rotate: [0, -10, 6, 0] }, { duration: τmax });
}
