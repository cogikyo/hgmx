import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { scale: [1, 0.9, 1] }, { duration: τmid });
  const thumbs = el.querySelectorAll("path");
  animate(thumbs, { translateY: [0, -2, 0.5, 0] }, { duration: 0.3, delay: stagger(0.09) });
}
