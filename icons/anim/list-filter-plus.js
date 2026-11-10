import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = ["M12 5H2", "M6 12h12", "M9 19h6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleX: [1, 0.5, 1.1, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const plus = el.querySelectorAll('path[d="M16 5h6"], path[d="M19 8V2"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
