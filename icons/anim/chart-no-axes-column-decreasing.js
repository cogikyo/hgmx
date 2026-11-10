import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = ["M5 21V3", "M12 21V9", "M19 21v-6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [1, 0.5, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
