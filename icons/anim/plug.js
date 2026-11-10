import { animate, τmax } from "../motion.js";

export default function (el) {
  const plug = el.querySelectorAll("path");
  animate(plug, { translateY: [0, -2, 1, 0] }, { duration: τmax });
}
