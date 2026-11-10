import { animate, τmax } from "../motion.js";

export default function (el) {
  const handle = el.querySelector("g");
  animate(handle, { scaleX: [1, 0.3, 1.1, 1] }, { duration: τmax });
}
