import { animate, τmax } from "../motion.js";

export default function (el) {
  const vest = el.querySelector("path");
  animate(vest, { translateY: [0, 2.5, -1, 0] }, { duration: τmax });
}
