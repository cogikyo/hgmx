import { animate, τmax } from "../motion.js";

export default function (el) {
  const engine = el.querySelectorAll("path");
  animate(engine, { translateX: [0, -2, 2, -1, 0] }, { duration: τmax });
}
