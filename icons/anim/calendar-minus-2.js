import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const minus = el.querySelector('path[d="M10 15h4"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
