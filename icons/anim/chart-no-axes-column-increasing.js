import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = ["M5 21v-6", "M12 21V9", "M19 21V3"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [1, 1.15, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
