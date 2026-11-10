import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lights = ["M7 13", "M12 14", "M17 13"].map((d) => el.querySelector(`path[d^="${d}h"]`));
  animate(lights, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const beam = el.querySelector("g");
  animate(beam, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
