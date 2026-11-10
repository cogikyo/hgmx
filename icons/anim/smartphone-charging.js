import { animate, τmax } from "../motion.js";

export default function (el) {
  const bolt = el.querySelector("path");
  animate(bolt, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmax });
}
