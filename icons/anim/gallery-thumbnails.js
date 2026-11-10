import { animate, stagger } from "../motion.js";

export default function (el) {
  const thumbs = el.querySelectorAll("path");
  animate(thumbs, { translateY: [0, -1.5, 0.5, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
