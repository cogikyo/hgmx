import { animate, τmax } from "../motion.js";

export default function (el) {
  const berry = el.querySelector("g");
  animate(berry, { scale: [1, 0.88, 1.06, 1] }, { duration: τmax });
}
