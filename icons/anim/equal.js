import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = [el.querySelector('line[y1="9"]'), el.querySelector('line[y1="15"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
