import { animate, τmid } from "../motion.js";

export default function (el) {
  const rocker = el.querySelector('path[d="M10 12h4"]');
  animate(rocker, { translateY: [0, 2, 0] }, { duration: τmid });
}
