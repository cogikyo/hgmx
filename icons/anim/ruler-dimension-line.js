import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ruler = el.querySelector("rect");
  animate(ruler, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ticks = ["M6", "M10", "M14", "M18"].map((d) => el.querySelector(`path[d="${d} 15v-3"]`));
  animate(ticks, { scale: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const span = el.querySelector('path[d="M22 6H2"]');
  animate(span, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
