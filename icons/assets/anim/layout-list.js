import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const thumbs = el.querySelectorAll("rect");
  animate(thumbs, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });
  const lines = el.querySelectorAll("path");
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
