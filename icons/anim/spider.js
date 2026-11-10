import { animate, τmax } from "../motion.js";

export default function (el) {
  const spider = el.querySelectorAll("path");
  animate(spider, { translateY: [-4, 1, 0] }, { duration: τmax });
}
